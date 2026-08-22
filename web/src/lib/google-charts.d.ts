declare const google: {
  charts: {
    load: (version: string, options: { packages: Array<string> }) => void
    setOnLoadCallback: (callback: () => void) => void
  }
  visualization: {
    arrayToDataTable: (data: Array<Array<unknown>>) => unknown
    PieChart: new (element: Element) => {
      draw: (data: unknown, options: Record<string, unknown>) => void
    }
  }
}
