package com.example.app

import org.springframework.boot.autoconfigure.SpringBootApplication
import org.springframework.boot.runApplication

@SpringBootApplication
class ExampleAppApplication

fun main(args: Array<String>) {
	runApplication<ExampleAppApplication>(*args)
}
