#!/usr/bin/env python3

import http.server
import socketserver
import sys
from pathlib import Path


class RepositoryHandler(http.server.BaseHTTPRequestHandler):
    release = ""

    def do_GET(self):
        body = (
            "<?xml version=\"1.0\" encoding=\"UTF-8\"?>"
            "<metadata><versioning><latest>{0}</latest><release>{0}</release>"
            "<versions><version>{0}</version></versions>"
            "<lastUpdated>20260824000000</lastUpdated></versioning></metadata>"
        ).format(self.release).encode("utf-8")
        self.send_response(200)
        self.send_header("Content-Type", "application/xml")
        self.send_header("Content-Length", str(len(body)))
        self.end_headers()
        self.wfile.write(body)

    def log_message(self, _format, *_args):
        return


class LoopbackServer(http.server.ThreadingHTTPServer):
    def server_bind(self):
        socketserver.TCPServer.server_bind(self)
        self.server_name = "localhost"
        self.server_port = self.server_address[1]


def main():
    if len(sys.argv) != 3:
        raise SystemExit("usage: loopback_repository.py PORT_FILE RELEASE")
    port_file, RepositoryHandler.release = sys.argv[1:]
    server = LoopbackServer(("127.0.0.1", 0), RepositoryHandler)
    Path(port_file).write_text(
        "http://127.0.0.1:{}\n".format(server.server_address[1]), encoding="utf-8"
    )
    server.serve_forever()


if __name__ == "__main__":
    main()
