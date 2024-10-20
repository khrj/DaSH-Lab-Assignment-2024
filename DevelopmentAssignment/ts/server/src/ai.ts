import { GoogleGenerativeAI } from "@google/generative-ai"

const apiKey = process.env.GEMINI_API_KEY
if (!apiKey) throw "Missing GEMINI_API_KEY in environment"

const genAI = new GoogleGenerativeAI(apiKey)
const model = genAI.getGenerativeModel({ model: "gemini-1.5-flash" })

const generationConfig = {
	temperature: 1,
	topP: 0.95,
	topK: 64,
	maxOutputTokens: 8192,
	responseMimeType: "text/plain",
}

export const ask = async (q: string): Promise<string> => {
	console.info(`⏳ Processing "${q}"...`)

	const c = model.startChat({
		generationConfig,
		history: [],
	})

	const result = await c.sendMessage(q)
	const content = result.response.text()

	console.info(`✅ Processed "${q}"...`)

	return content
}
