import { API } from "@/lib/api.consts"
import { captures } from "@/lib/captures.actions"
import { customStore } from "@/lib/custom-ingredients.store"
import { createSignal, onMount } from "solid-js"

export default function DishFetch() {
	onMount(() => {
		const ingredients = (window.history.state || new URLSearchParams(window.location.search).getAll('ingredient')) as string[]
		const [_, set] = customStore

		set(ingredients.map(i => createSignal(i)))
	})

	const onClick = () => {
		const [customs] = customStore
		const flatCustoms = customs().map(c => {
			const [ingredients] = c
			return ingredients()
		})

		const ingredients = captures().flatMap(c => {
			const [ingredients] = c.ingredients
			return ingredients()
		})
			.concat(flatCustoms)
			.map(i => i.trim())
			.filter(Boolean)

		const params = new URLSearchParams()
		ingredients.forEach(ingredient => {
			params.append('ingredient', ingredient)
		})

		const query = '?' + params.toString()
		window.history.pushState(ingredients, '', query)
		fetch(`${API}/dishes${query}`).then(response => response.json()).then((data) => {
			console.log(data)
		})
	}

	return (
		<div class="grid place-content-center">
			<button class="bg-sky-500 hover:bg-sky-700 text-white font-bold py-2 px-4 rounded" onClick={onClick}>Search dishes</button>
		</div>
	)
}
