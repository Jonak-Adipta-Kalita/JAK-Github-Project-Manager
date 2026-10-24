<script lang="ts">
    import { SyncLoader } from "svelte-loading-spinners";
    import axios from "axios";
    import { z } from "zod";

    let firstTerm = "0";
    let commonDifference = "0";
    let numberOfTerms = "1";

    let display = false;
    let errors: {
        firstTerm?: string[] | undefined;
        commonDifference?: string[] | undefined;
        numberOfTerms?: string[] | undefined;
    } | null = null;
    let disabledButton = false;

    let series: null | string = null;
    let sum: null | number = null;

    const zodNumericString = (terms?: boolean) => {
        return z
            .string()
            .regex(terms ? /^\d+$/ : /-?\d+(\.\d+)?$/, {
                message: "Invalid number"
            })
            .transform(Number);
    };

    const formSchema = z.object({
        firstTerm: zodNumericString(),
        commonDifference: zodNumericString(),
        numberOfTerms: zodNumericString(true)
    });

    const piston = async (need: "series" | "sum") => {
        const sumSourceCode = await (await fetch("/sum.py")).text();
        const seriesSourceCode = await (await fetch("/series.py")).text();

        const files = [
            {
                name: "sum.py",
                content: sumSourceCode
            },
            {
                name: "series.py",
                content: seriesSourceCode
            }
        ];

        return (
            await axios.post(
                "https://emkc.org/api/v2/piston/execute",
                {
                    language: "python",
                    version: "3.10.0",
                    files: need === "series" ? files.reverse() : files,
                    args: [firstTerm, commonDifference, numberOfTerms]
                },
                {
                    headers: {
                        "Content-Type": "application/json"
                    }
                }
            )
        ).data.run.output;
    };

    const calculate = async (
        e: SubmitEvent & {
            currentTarget: EventTarget & HTMLFormElement;
        }
    ) => {
        const parseResult = await formSchema.safeParseAsync({
            firstTerm,
            commonDifference,
            numberOfTerms
        });

        if (!parseResult.success) {
            errors = parseResult.error.flatten().fieldErrors;
            return;
        } else {
            errors = null;
        }

        display = true;

        series = null;
        sum = null;

        disabledButton = true;

        series = ((await piston("series")) as string)
            .split(", ")
            .map((num) => Number(num))
            .join(", ");
        sum = Number(await piston("sum"));

        disabledButton = false;
    };
</script>

<svelte:head>
    <title>JAK Arithmatic Progression App</title>
    <meta name="description" content="For Priyanka Ma'am 💙" />
</svelte:head>

<form
    on:submit|preventDefault={calculate}
    class="flex flex-col items-center justify-center bg-gray-900 shadow-lg h-screen space-y-6"
>
    <div
        class="flex flex-col items-center justify-center bg-gray-800 rounded-lg space-y-4 px-10 py-5"
    >
        <div class="inputContainer">
            <p class="inputLabels">First Term:</p>
            <ul class="text-red-500">
                {#if errors?.firstTerm}
                    {#each errors?.firstTerm as error}
                        <li>{error}</li>
                    {/each}
                {/if}
            </ul>
            <input
                class="inputs"
                placeholder="First Term"
                disabled={disabledButton}
                bind:value={firstTerm}
                required
            />
        </div>
        <div class="inputContainer">
            <p class="inputLabels">Common Difference:</p>
            <ul class="text-red-500">
                {#if errors?.commonDifference}
                    {#each errors?.commonDifference as error}
                        <li>{error}</li>
                    {/each}
                {/if}
            </ul>
            <input
                class="inputs"
                placeholder="Common Difference"
                disabled={disabledButton}
                bind:value={commonDifference}
                required
            />
        </div>
        <div class="inputContainer">
            <p class="inputLabels">No. of Terms:</p>
            <ul class="text-red-500">
                {#if errors?.numberOfTerms}
                    {#each errors?.numberOfTerms as error}
                        <li>{error}</li>
                    {/each}
                {/if}
            </ul>
            <input
                class="inputs"
                placeholder="No. of Terms"
                disabled={disabledButton}
                bind:value={numberOfTerms}
                required
                min="1"
            />
        </div>
        {#if display}
            {#if series !== null && sum !== null}
                <div class="mt-10" />
                <textarea
                    class="bg-gray-300 w-64 px-4 py-2 rounded-lg text-sm text-gray-900 font-bold outline-none text-justify"
                    value={`AP: ${series}\nSum: ${sum}`}
                    readonly
                    rows={5}
                />
            {:else}
                <SyncLoader size="60" color="#d1d5db" unit="px" duration="1s" />
            {/if}
        {/if}
    </div>
    <button
        type="submit"
        class="w-80 py-[10px] rounded-lg bg-blue-500 hover:bg-blue-600 focus:outline-none focus:ring-2 focus:ring-blue-500 focus:ring-offset-2 px-10"
    >
        <span class="text-xl font-semibold text-white">Calculate</span>
    </button>
</form>
