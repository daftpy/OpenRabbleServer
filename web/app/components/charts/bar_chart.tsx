import { Bar } from "react-chartjs-2";
import type { ChartData, ChartOptions } from "chart.js";

export type BarData = ChartData<"bar", number[], string>;

export type AspectRatio = number | [number, number];
const ratioToString = (aspectRatio : AspectRatio) => {
  if (typeof aspectRatio == "number") {
    return `${aspectRatio}`
  } else if (Array.isArray(aspectRatio)) {
    return `${aspectRatio[0]}:${aspectRatio[1]}`
  }
  return "2" // default fallback
}

export function BarChart({ data, options } : { data: BarData, options?: ChartOptions<"bar">}) {

  return (
    <Bar data={data} options={options} />
  )
}