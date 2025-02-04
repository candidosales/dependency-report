import { AfterViewInit, Component, ElementRef, input, viewChild } from '@angular/core';
declare var google: any;

@Component({
    selector: 'app-treemap-chart',
    templateUrl: './treemap-chart.component.html',
    standalone: true,
})
export class TreemapChartComponent implements AfterViewInit{

  treemapChart = viewChild<ElementRef>('treemapChart');
  data = input<any>();
  title = input('');
  height = input('250px');
  width = input('100%');


  drawChart = () => {
    const data = google.visualization.arrayToDataTable(this.data());

    const options = {
        highlightOnMouseOver: true,
        maxDepth: 1,
        maxPostDepth: 2,
        minColor: '#7cd6fd',
        midColor: '#5e64ff',
        maxColor: '#ffa00a',
        headerHeight: 30,
        fontColor: 'black',
        showScale: true,
        height: 500,
        useWeightedAverageForAggregation: true
    };

    const chart = new google.visualization.TreeMap(this.treemapChart().nativeElement);
    // const chart = new google.visualization.TreeMap(document.getElementById('chart_div'));


    chart.draw(data, options);
  }

  ngAfterViewInit() {
    google.charts.load('current', { packages: ['treemap'] });
    google.charts.setOnLoadCallback(this.drawChart);
  }
}