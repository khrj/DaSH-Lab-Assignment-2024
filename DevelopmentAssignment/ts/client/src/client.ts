import { io } from "socket.io-client"
import { writeFile, readFile } from "fs/promises"

const socket = io("http://server:3000")

let responses = []
let questions: string[] = []

socket.on("response", (q, r) => {
	const output = {
		clientId: socket.id,
		prompt: q.prompt,
		message: r,
		timeSent: q.timeSent,
		timeRecvd: q.timeRecvd,
		source: questions.includes(q.prompt) ? "gemini" : "user",
	}

	responses.push(output)

	if (responses.length >= 12) {
		writeFile(`outputs/output-${socket.id}.json`, JSON.stringify(responses))
			.then(() => {
				console.log("Received all responses, exiting...")
				process.exit(0)
			})
			.catch(console.error)
	}
})

const sleep = (ms: number) => new Promise(r => setTimeout(r, ms))

socket.on("connect", async () => {
	console.log("connected")
	const txt = await readFile("input.txt", { encoding: "utf8" })

	questions = txt.split("\n")

	for (const q of questions) {
		await sleep(4000) // avoid gemini rate limits
		socket.emit("query", {
			prompt: q,
			timeSent: Date.now(),
		})
	}
})
