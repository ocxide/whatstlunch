import { insertCapture, removeCapture } from "@/lib/captures.actions"
import { For, createSignal, onCleanup } from "solid-js"

type FileRead = {
	file: File,
	blob: string
}

export default function ImageCaptures() {
	let addInput: HTMLInputElement
	const [captures, setCaptures] = createSignal<FileRead[]>([])

	const add = (file: File) => {
		const blob = URL.createObjectURL(file)
		setCaptures(captures => [...captures, { file, blob }])
		addInput.value = ''

		insertCapture(file)
	}

	const updateOne = (file: File, i: number) => {
		const blob = URL.createObjectURL(file)

		const previousFilename = captures()[i].file.name

		setCaptures(captures => {
			captures[i] = { file, blob }
			return captures.slice()
		})

		insertCapture(file, previousFilename)
	}

	const removeOne = (i: number) => {
		const capture = captures()[i]

		removeCapture(capture.file.name)

		setCaptures(captures => {
			captures.splice(i, 1).forEach(read => URL.revokeObjectURL(read.blob))
			return captures.slice()
		})
	}

	onCleanup(() => {
		captures().forEach(read => {
			URL.revokeObjectURL(read.blob)
		})
	})

	return <div class="grid gap-2">
		<ul class="grid gap-4">
			<For each={captures()}>{
				(read, i) => (<li>
					<p>{read.file.name}</p>
					<img src={read.blob} alt="" />
					<div class="flex gap-x-2">
						<input type="file" accept="image/*;capture=camera" onInput={e => updateOne(e.target.files![0], i())} />
						<button onClick={() => removeOne(i())}>Remove</button>
					</div>
				</li>)}
			</For>
		</ul>

		<hr />

		<input ref={e => addInput = e} type="file" accept="image/*;capture=camera" onInput={e => add(e.target.files![0])} />
	</div >
}
