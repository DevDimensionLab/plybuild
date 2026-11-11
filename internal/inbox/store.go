package inbox

import (
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

const (
	messagesName       = "messages"
	readName           = "read"
	conversationLength = len("conv_") + 32
)

var (
	messageNamePattern = regexp.MustCompile(`^([0-9]{8}T[0-9]{6}\.[0-9]{9}Z)_(conv_[0-9a-f]{32})_(msg_[0-9a-f]{32})\.json$`)
	markNamePattern    = regexp.MustCompile(`^(conv_[0-9a-f]{32})_(msg_[0-9a-f]{32})\.json$`)
	errNotRealDir      = itemInvalid("store path component is not a real directory")
)

// Put stores the reply once. A replay of the same payload is unchanged and
// keeps everything stored; a different payload under the same identity is a
// conflict. Both created and unchanged confirm that the payload is stored.
func (s *Store) Put(req PutRequest) (PutResult, error) {
	message, err := s.newMessage(req)
	if err != nil {
		return PutResult{}, err
	}
	if err := os.MkdirAll(s.root, dirMode); err != nil {
		return PutResult{}, ioError("create inbox directory", err)
	}
	unlock, err := s.lock()
	if err != nil {
		return PutResult{}, err
	}
	defer unlock()
	for _, dir := range []string{s.ownerDir(), s.messagesDir()} {
		if err := ensureDirectory(dir); err != nil {
			return PutResult{}, err
		}
	}
	existing, found, err := s.lookup(message.ConversationID, message.MessageID)
	if err != nil {
		return PutResult{}, err
	}
	if found {
		if existing.PayloadSHA256 == message.PayloadSHA256 {
			return PutResult{PutResultKind, SchemaVersion, EffectUnchanged, existing}, nil
		}
		return PutResult{}, &Error{Class: ClassConflict, Detail: "payload differs", ConflictingFields: conflictingFields(payloadOf(existing), payloadOf(message)),
			ExistingPayloadSHA256: existing.PayloadSHA256, SubmittedPayloadSHA256: message.PayloadSHA256}
	}
	contents, err := encodeJSON(message)
	if err != nil {
		return PutResult{}, ioError("encode message", err)
	}
	if err := s.replace(filepath.Join(s.messagesDir(), messageKey(message)+fileSuffix), contents, fileMode); err != nil {
		return PutResult{}, ioError("write message", err)
	}
	return PutResult{PutResultKind, SchemaVersion, EffectCreated, message}, nil
}

// newMessage validates req and builds the message a first write stores.
func (s *Store) newMessage(req PutRequest) (Message, error) {
	if req.Kind != PutRequestKind {
		return Message{}, invalidArguments("kind must be " + PutRequestKind)
	}
	p, err := checkPayload(req)
	if err != nil {
		return Message{}, err
	}
	if err := checkMetadata(req.AgentLabel, req.Provider, req.WriteKind); err != nil {
		return Message{}, err
	}
	hash, err := payloadSHA256(s.project, s.owner, p)
	if err != nil {
		return Message{}, err
	}
	return Message{
		Kind: MessageKind, SchemaVersion: SchemaVersion, ProjectID: s.project, Owner: s.owner,
		ConversationID: p.conversationID, MessageID: p.messageID, RepliedAtUTC: p.repliedAtUTC,
		Title: p.title, DerivedTitle: DeriveTitle(p.body), DerivedTitleRule: DerivedTitleRule,
		Body: p.body, BodySHA256: bodySHA256(p.body), Truncated: p.truncated, FullSHA256: p.fullSHA256,
		AnswerSource: p.answerSource, AgentLabel: req.AgentLabel, Provider: req.Provider, WriteKind: req.WriteKind,
		RecordedAtUTC: s.now().UTC().Format(canonicalLayout), PayloadSHA256: hash,
	}, nil
}

// Show returns one item with its read state.
func (s *Store) Show(conversation, message string) (Item, error) {
	if err := checkIdentity(conversation, message); err != nil {
		return Item{}, err
	}
	stored, found, err := s.lookup(conversation, message)
	if err != nil {
		return Item{}, err
	}
	if !found {
		return Item{}, notFound(conversation, message)
	}
	return newItem(stored, s.markState(conversation, message)), nil
}

// Mark sets or removes the read mark of an existing item. It never touches the message.
func (s *Store) Mark(conversation, message string, value MarkValue) (MarkResult, error) {
	if err := checkIdentity(conversation, message); err != nil {
		return MarkResult{}, err
	}
	if value != MarkRead && value != MarkUnread {
		return MarkResult{}, invalidArguments("value must be read or unread")
	}
	if info, err := os.Stat(s.root); errors.Is(err, fs.ErrNotExist) {
		return MarkResult{}, notFound(conversation, message)
	} else if err != nil {
		return MarkResult{}, ioError("check inbox directory", err)
	} else if !info.IsDir() {
		return MarkResult{}, errNotRealDir
	}
	unlock, err := s.lock()
	if err != nil {
		return MarkResult{}, err
	}
	defer unlock()
	if _, found, err := s.lookup(conversation, message); err != nil {
		return MarkResult{}, err
	} else if !found {
		return MarkResult{}, notFound(conversation, message)
	}
	result := MarkResult{MarkResultKind, SchemaVersion, conversation, message, value == MarkRead, EffectNone}
	path := s.markPath(conversation, message)
	present := false
	if exists, err := realDirectory(s.readDir()); err != nil {
		return MarkResult{}, err
	} else if exists {
		contents, err := readStrict(path)
		present = err == nil
		switch {
		case errors.Is(err, fs.ErrNotExist):
		case err != nil:
			return MarkResult{}, err
		default:
			if _, err := decodeMarkFile(contents, filepath.Base(path), s.owner); err != nil {
				return MarkResult{}, itemInvalid("read mark is invalid")
			}
		}
	}
	switch {
	case value == MarkRead && !present:
		if err := ensureDirectory(s.readDir()); err != nil {
			return MarkResult{}, err
		}
		mark := ReadMark{ReadMarkKind, SchemaVersion, s.owner, conversation, message, s.now().UTC().Format(canonicalLayout)}
		encoded, err := encodeJSON(mark)
		if err != nil {
			return MarkResult{}, ioError("encode read mark", err)
		}
		if err := s.replace(path, encoded, fileMode); err != nil {
			return MarkResult{}, ioError("write read mark", err)
		}
		result.Effect = EffectMarked
	case value == MarkUnread && present:
		if err := os.Remove(path); err != nil {
			return MarkResult{}, ioError("remove read mark", err)
		}
		if err := syncDirectory(s.readDir()); err != nil {
			return MarkResult{}, ioError("sync read directory", err)
		}
		result.Effect = EffectMarked
	}
	return result, nil
}

// List returns one page in (replied_at_utc desc, conversation_id, message_id)
// order after cursor. Counts come from file names; only page items are decoded.
func (s *Store) List(limit int, cursor string) (List, error) {
	if limit < 1 || limit > 50 {
		return List{}, invalidArguments("--limit must be 1-50")
	}
	result := List{Kind: ListKind, SchemaVersion: SchemaVersion, ProjectID: s.project, Owner: s.owner, Limit: limit,
		ReadState: ReadStateAvailable, Items: []Item{}, Invalid: []InvalidEntry{}}
	after := ""
	if cursor != "" {
		key, err := s.decodeCursor(cursor)
		if err != nil {
			return List{}, err
		}
		after, result.Cursor = key, &cursor
	}
	keys, err := s.messageKeys(&result)
	if err != nil {
		return List{}, err
	}
	read, unknown, available := s.marks(keys, &result)
	sort.Slice(keys, func(i, j int) bool { return keyLess(keys[i], keys[j]) })
	result.TotalCount = len(keys)
	start := sort.Search(len(keys), func(i int) bool { return after == "" || keyLess(after, keys[i]) })
	end := start + limit
	if end > len(keys) {
		end = len(keys)
	}
	pageValid := true
	for _, key := range keys[start:end] {
		conversation, message := keyIdentity(key)
		stored, err := s.readMessage(key + fileSuffix)
		if err != nil {
			pageValid = false
			result.Invalid = append(result.Invalid, InvalidEntry{key + fileSuffix, errorText(err)})
			continue
		}
		var state *bool
		if available && !unknown[identity(conversation, message)] {
			isRead := read[identity(conversation, message)]
			state = &isRead
		}
		result.Items = append(result.Items, newItem(stored, state))
	}
	if end < len(keys) {
		next := s.encodeCursor(keys[end-1])
		result.NextCursor = &next
	}
	if !available {
		result.UnreadCountBasis = BasisUnavailable
		return result, nil
	}
	unread := 0
	for _, key := range keys {
		if id := identity(keyIdentity(key)); !read[id] && !unknown[id] {
			unread++
		}
	}
	result.UnreadCount = &unread
	result.UnreadCountExact = pageValid && start == 0 && end == len(keys) && len(result.Invalid) == 0
	result.UnreadCountBasis = BasisFilenameUnverified
	if result.UnreadCountExact {
		result.UnreadCountBasis = BasisValidatedItems
	}
	return result, nil
}

// messageKeys returns the keys of valid, unique message names and reports every other name.
func (s *Store) messageKeys(result *List) ([]string, error) {
	for _, dir := range []string{s.ownerDir(), s.messagesDir()} {
		if exists, err := realDirectory(dir); err != nil || !exists {
			return []string{}, err
		}
	}
	entries, err := os.ReadDir(s.messagesDir())
	if err != nil {
		return nil, ioError("read messages directory", err)
	}
	byIdentity := map[string][]string{}
	order := []string{}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		match := messageNamePattern.FindStringSubmatch(name)
		switch {
		case match == nil || !validCompactTime(match[1]):
			result.Invalid = append(result.Invalid, InvalidEntry{name, ClassItemInvalid + ": invalid file name"})
		case !entry.Type().IsRegular():
			result.Invalid = append(result.Invalid, InvalidEntry{name, ClassItemInvalid + ": not a regular file"})
		default:
			id := identity(match[2], match[3])
			if byIdentity[id] == nil {
				order = append(order, id)
			}
			byIdentity[id] = append(byIdentity[id], strings.TrimSuffix(name, fileSuffix))
		}
	}
	keys := []string{}
	for _, id := range order {
		if names := byIdentity[id]; len(names) > 1 {
			for _, key := range names {
				result.Invalid = append(result.Invalid, InvalidEntry{key + fileSuffix, ClassItemInvalid + ": duplicate identity"})
			}
			continue
		}
		keys = append(keys, byIdentity[id][0])
	}
	return keys, nil
}

// marks validates every read mark. A valid mark means read, an invalid one
// makes that item's state unknown; available is false when read/ is unreadable.
func (s *Store) marks(keys []string, result *List) (read, unknown map[string]bool, available bool) {
	read, unknown = map[string]bool{}, map[string]bool{}
	exists, err := realDirectory(s.readDir())
	if err != nil {
		result.ReadState = ReadStateUnavailable
		return read, unknown, false
	}
	if !exists {
		return read, unknown, true
	}
	entries, err := os.ReadDir(s.readDir())
	if err != nil {
		result.ReadState = ReadStateUnavailable
		return read, unknown, false
	}
	items := map[string]bool{}
	for _, key := range keys {
		items[identity(keyIdentity(key))] = true
	}
	for _, entry := range entries {
		name := entry.Name()
		if strings.HasPrefix(name, ".") {
			continue
		}
		match := markNamePattern.FindStringSubmatch(name)
		if match == nil {
			result.Invalid = append(result.Invalid, InvalidEntry{name, ClassItemInvalid + ": invalid file name"})
			continue
		}
		id := identity(match[1], match[2])
		contents, err := readStrict(filepath.Join(s.readDir(), name))
		if err == nil {
			_, err = decodeMarkFile(contents, name, s.owner)
		}
		switch {
		case err != nil:
			unknown[id] = true
			result.Invalid = append(result.Invalid, InvalidEntry{name, errorText(err)})
		case !items[id]:
			result.Invalid = append(result.Invalid, InvalidEntry{name, ClassItemInvalid + ": read mark without message"})
		default:
			read[id] = true
		}
	}
	return read, unknown, true
}

// markState is true, false, or nil when the mark state is unreadable.
func (s *Store) markState(conversation, message string) *bool {
	isRead := false
	if exists, err := realDirectory(s.readDir()); err != nil {
		return nil
	} else if !exists {
		return &isRead
	}
	path := s.markPath(conversation, message)
	contents, err := readStrict(path)
	if errors.Is(err, fs.ErrNotExist) {
		return &isRead
	}
	if err != nil {
		return nil
	}
	if _, err := decodeMarkFile(contents, filepath.Base(path), s.owner); err != nil {
		return nil
	}
	isRead = true
	return &isRead
}

// lookup finds the message with this identity. Missing directories mean not
// found; two names with the identity are invalid.
func (s *Store) lookup(conversation, message string) (Message, bool, error) {
	for _, dir := range []string{s.ownerDir(), s.messagesDir()} {
		if exists, err := realDirectory(dir); err != nil || !exists {
			return Message{}, false, err
		}
	}
	entries, err := os.ReadDir(s.messagesDir())
	if err != nil {
		return Message{}, false, ioError("read messages directory", err)
	}
	suffix := "_" + identity(conversation, message) + fileSuffix
	matches := []string{}
	for _, entry := range entries {
		if name := entry.Name(); !strings.HasPrefix(name, ".") && strings.HasSuffix(name, suffix) {
			matches = append(matches, name)
		}
	}
	switch len(matches) {
	case 0:
		return Message{}, false, nil
	case 1:
		stored, err := s.readMessage(matches[0])
		return stored, err == nil, err
	}
	return Message{}, false, itemInvalid("duplicate identity")
}

// readMessage treats an unreadable message file as invalid, never as retryable (§3.3).
func (s *Store) readMessage(name string) (Message, error) {
	contents, err := readStrict(filepath.Join(s.messagesDir(), name))
	if errors.Is(err, fs.ErrNotExist) {
		return Message{}, itemInvalid("message file vanished")
	}
	var inboxError *Error
	if errors.As(err, &inboxError) && inboxError.Class == ClassIOError {
		return Message{}, itemInvalid("message file is unreadable: " + inboxError.Detail)
	}
	if err != nil {
		return Message{}, err
	}
	return decodeMessageFile(contents, name, s.project, s.owner)
}

func (s *Store) lock() (func(), error) {
	file, err := openLockFile(filepath.Join(s.root, lockName))
	if err != nil {
		return nil, ioError("open inbox lock", err)
	}
	deadline := time.Now().Add(s.lockTimeout)
	for {
		locked, err := tryLock(file)
		if err != nil {
			_ = file.Close()
			return nil, ioError("lock inbox", err)
		}
		if locked {
			return func() { _ = unlockFile(file); _ = file.Close() }, nil
		}
		if time.Now().After(deadline) {
			_ = file.Close()
			return nil, &Error{Class: ClassBusy, Detail: "inbox lock not acquired within 2 s"}
		}
		time.Sleep(10 * time.Millisecond)
	}
}

func (s *Store) ownerDir() string    { return filepath.Join(s.root, s.owner) }
func (s *Store) messagesDir() string { return filepath.Join(s.ownerDir(), messagesName) }
func (s *Store) readDir() string     { return filepath.Join(s.ownerDir(), readName) }
func (s *Store) markPath(conversation, message string) string {
	return filepath.Join(s.readDir(), identity(conversation, message)+fileSuffix)
}

func identity(conversation, message string) string { return conversation + "_" + message }

func messageKey(m Message) string {
	return compactTime(m.RepliedAtUTC) + "_" + identity(m.ConversationID, m.MessageID)
}

// keyIdentity splits a valid key <T>_<conversation>_<message>; all three parts are fixed width.
func keyIdentity(key string) (string, string) {
	conversation := key[len(compactLayout)+1 : len(compactLayout)+1+conversationLength]
	return conversation, key[len(compactLayout)+1+conversationLength+1:]
}

// keyLess orders newest first, then by conversation and message ascending.
// The time part is fixed width, so string comparison is time order.
func keyLess(a, b string) bool {
	at, bt := a[:len(compactLayout)], b[:len(compactLayout)]
	if at != bt {
		return at > bt
	}
	return a[len(at):] < b[len(bt):]
}

func validCompactTime(value string) bool {
	parsed, err := time.Parse(compactLayout, value)
	return err == nil && parsed.Format(compactLayout) == value
}

func checkIdentity(conversation, message string) error {
	if !conversationPattern.MatchString(conversation) || !messagePattern.MatchString(message) {
		return invalidArguments("conversation and message must be conv_ and msg_ IDs with 32 lowercase hex")
	}
	return nil
}

// realDirectory reports whether path exists as a directory that is not a symlink.
func realDirectory(path string) (bool, error) {
	info, err := os.Lstat(path)
	switch {
	case errors.Is(err, fs.ErrNotExist):
		return false, nil
	case err != nil:
		return false, ioError("check store directory", err)
	case !info.IsDir() || info.Mode()&fs.ModeSymlink != 0:
		return false, errNotRealDir
	}
	return true, nil
}

func ensureDirectory(path string) error {
	if err := os.Mkdir(path, dirMode); err != nil && !errors.Is(err, fs.ErrExist) {
		return ioError("create store directory", err)
	}
	if exists, err := realDirectory(path); err != nil {
		return err
	} else if !exists {
		return ioError("create store directory", fs.ErrNotExist)
	}
	return nil
}

// readStrict reads a regular file of at most 1 MiB without following a final symlink.
// A missing file returns an error matching fs.ErrNotExist.
func readStrict(path string) ([]byte, error) {
	file, err := openNoFollow(path)
	if errors.Is(err, fs.ErrNotExist) {
		return nil, fs.ErrNotExist
	}
	if errors.Is(err, errNotRegular) {
		return nil, itemInvalid("not a regular file")
	}
	if err != nil {
		return nil, ioError("open store file", err)
	}
	defer file.Close()
	info, err := file.Stat()
	if err != nil {
		return nil, ioError("stat store file", err)
	}
	if !info.Mode().IsRegular() {
		return nil, itemInvalid("not a regular file")
	}
	contents, err := io.ReadAll(io.LimitReader(file, maxFileBytes+1))
	if err != nil {
		return nil, ioError("read store file", err)
	}
	if len(contents) > maxFileBytes {
		return nil, itemInvalid("file exceeds 1 MiB")
	}
	return contents, nil
}

var errNotRegular = errors.New("not a regular file")

func errorText(err error) string {
	var inboxError *Error
	if errors.As(err, &inboxError) {
		return inboxError.Error()
	}
	return ClassIOError + ": unexpected error"
}
