import { AfterViewInit, Component, ElementRef, ViewChild, Input } from '@angular/core';
declare var google: any;

@Component({
  selector: 'app-treemap-chart',
  templateUrl: './treemap-chart.component.html',
})
export class TreemapChartComponent implements AfterViewInit{

  @ViewChild('treemapChart') treemapChart: ElementRef;
  @Input() data: any;
  @Input() title = '';
  @Input() height = '250px';
  @Input() width = '100%';


  drawChart = () => {
    const data = google.visualization.arrayToDataTable(this.data);

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

    const chart = new google.visualization.TreeMap(this.treemapChart.nativeElement);
    // const chart = new google.visualization.TreeMap(document.getElementById('chart_div'));


    chart.draw(data, options);
  }

  ngAfterViewInit() {
    google.charts.load('current', { packages: ['treemap'] });
    google.charts.setOnLoadCallback(this.drawChart);
  }
}