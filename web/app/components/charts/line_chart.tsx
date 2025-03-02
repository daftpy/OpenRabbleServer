import type { ChartData, ChartOptions, Point } from "chart.js";
import { Line } from "react-chartjs-2";

export type LineData = ChartData<"line", (number | Point | null)[], unknown>

export function LineChart({ data, options } : { data: LineData, options: ChartOptions<"line">}) {
  return (
    <Line data={data} options={options} />
  )
}