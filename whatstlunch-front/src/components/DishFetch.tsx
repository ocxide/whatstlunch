import { captures } from "@/lib/captures.actions"
import { customStore } from "@/lib/custom-ingredients.store"

export default function DishFetch() {

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

		console.log(params.toString())
	}

	return (
		<button onClick={onClick}>Search dishes</button>
	)
}
