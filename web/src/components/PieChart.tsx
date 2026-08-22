import { useEffect, useRef } from 'react'

interface PieChartProps {
  data: Array<[string, unknown]> | undefined
  title?: string
  height?: string
  width?: string
}

const COLORS = [
  '#D53F8C',
  '#805AD5',
  '#5A67D8',
  '#3182CE',
  '#319795',
  '#38A169',
  '#D69E2E',
]

function getGoogleCharts() {
  return (window as unknown as { google?: typeof google }).google
}

export default function PieChart({
  data,
  title = '',
  height = '250px',
  width = '100%',
}: PieChartProps) {
  const containerRef = useRef<HTMLDivElement>(null)

  useEffect(() => {
    let cancelled = false

    function draw() {
      const g = getGoogleCharts()
      if (cancelled || !data || !g?.visualization || !containerRef.current) {
        return
      }

      try {
        const table = g.visualization.arrayToDataTable(data)
        const chart = new g.visualization.PieChart(containerRef.current)
        chart.draw(table, {
          title,
          legend: { position: 'right' },
          colors: COLORS,
        })
      } catch (err) {
        console.error(err)
      }
    }

    function loadAndDraw() {
      const g = getGoogleCharts()
      if (!g) {
        window.setTimeout(loadAndDraw, 100)
        return
      }
      g.charts.load('current', { packages: ['corechart'] })
      g.charts.setOnLoadCallback(draw)
    }

    loadAndDraw()
    window.addEventListener('resize', draw)

    return () => {
      cancelled = true
      window.removeEventListener('resize', draw)
    }
  }, [data, title])

  return <div ref={containerRef} style={{ height, width }} />
}
