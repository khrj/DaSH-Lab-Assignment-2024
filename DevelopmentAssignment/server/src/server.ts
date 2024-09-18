import "dotenv/config"
import { Server } from "socket.io"
import { ask } from "./ai.js"

const io = new Server()

io.on("connection", socket =>
	socket.on("query", q =>
		ask(q.prompt).then(resp =>
			io.emit(
				"response",
				{
					...q,
					timeRecvd: Date.now(),
				},
				resp,
			),
		),
	),
)

io.listen(3000)
console.log("Socket IO server listening on 3000")
