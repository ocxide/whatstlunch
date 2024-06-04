import { API } from "@/lib/api.consts"
import { captures } from "@/lib/captures.actions"
import { customStore } from "@/lib/custom-ingredients.store"
import { For, Show, createSignal, onMount, type Accessor, type Signal } from "solid-js"

export type Dish = {
	title: string;
	introduction: string;
	duration: string;
	foodType: string;
	ingredients: string[];
	preparation: string[];
}

export default function DishFetch() {
	const [dishes, setDishes] = createSignal<Dish[]>([])
	const [dishEntryOpened, setDishEntryOpened] = createSignal(-1)


	const percentageSignal = createSignal(0)
	const [percentage, setPercentage] = percentageSignal

	const [limit, setLimit] = createSignal(0)

	const [isPercentage, setIsPercentage] = createSignal(false)

	onMount(() => {
		const search = new URLSearchParams(window.location.search)
		const ingredients = search.getAll('ingredient')

		if (ingredients.length > 0) {
			const [_, set] = customStore
			set(ingredients.map(i => createSignal(i)))
		}

		const requireStr = search.get('require') ?? '0'
		const requireInt = parseInt(requireStr, 10)

		if (isNaN(requireInt)) return

		if (requireInt === 0 && requireStr.includes('.')) {
			setIsPercentage(true)
			setPercentage(Number(requireStr))

			return
		}

		setLimit(requireInt)
		setIsPercentage(false)
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

		params.append('require', (isPercentage() ? percentage() : limit()).toString())

		const query = '?' + params.toString()
		window.history.pushState(ingredients, '', query)

		fetch(`${API}/dishes${query}`).then(response => response.json()).then((data: Dish[]) => {
			setDishes(data)
		})
	}

	const onToggle = (index: number) => {
		setDishEntryOpened(i => i === index ? -1 : index)
	}

	return (
		<div class="">
			<div class="p-4 grid justify-center">
				<div class="min-h-10 flex-col flex justify-center items-center">
					<div class="grid sm:grid-cols-[auto,14rem,auto,auto] items-center gap-2">
						<span>Require at least:&nbsp;</span>

						<Show when={isPercentage()} fallback={<input type="number" value={limit()} onInput={e => setLimit(parseInt(e.currentTarget.value))} />}>
							<PercentageControl percentage={percentageSignal} />
						</Show>

						<input id="isPercentage" type="checkbox" checked={isPercentage()} onInput={e => setIsPercentage(e.currentTarget.checked)} />
						<label for="isPercentage">Is Percentage</label>
					</div>
				</div>

				<button class="bg-sky-500 hover:bg-sky-700 text-white font-bold py-2 px-4 rounded" onClick={onClick}>Search dishes</button>
			</div>

			<ul>
				<For each={dishes()}>
					{(dish, i) => (
						<DishEntry dish={dish} opened={() => dishEntryOpened() === i()} onToggle={() => onToggle(i())} />
					)}
				</For>
			</ul>
		</div>
	)
}

function PercentageControl({ percentage: signal }: { percentage: Signal<number> }) {
	const [percentage, setPercentage] = signal
	const percentageDisplay = () => Math.round(percentage() * 100)

	return (<div class="flex gap-x-2 items-center">
		<input
			type="range"
			class="min-w-0"
			min="0"
			max="1"
			step="0.05"
			value={percentage()}
			onInput={e => setPercentage(parseFloat(e.currentTarget.value))}
		/>
		<span>{percentageDisplay()}%</span>
	</div>)
}

function DishEntry({
	dish: { title, duration, foodType, ingredients, preparation, introduction },
	opened,
	onToggle
}: { dish: Dish, opened: Accessor<boolean>, onToggle: () => void }) {
	return (
		<li>
			<button onClick={onToggle} class="font-bold">{title}</button>

			<div class={`grid gap-2 ${!opened() ? 'hidden' : ''}`}>
				<p>{introduction}</p>
				<p>
					<span>Duration: {duration}</span>
					<span>&nbsp;-&nbsp;</span>
					<span>Food type: {foodType}</span>
				</p>

				<div>
					<p>Ingredients:</p>
					<ul class="list-disc list-inside">
						<For each={ingredients}>{ingredient => <li>{ingredient}</li>}</For>
					</ul>
				</div>

				<div>
					<p>Preparation:</p>
					<ul class="list-disc list-inside">
						<For each={preparation}>{step => <li>{step}</li>}</For>
					</ul>
				</div>
			</div>
		</li>
	)
}
