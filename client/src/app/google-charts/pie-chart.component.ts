import { AfterViewInit, Component, ElementRef, ViewChild, Input } from '@angular/core';
declare var google: any;

@Component({
  selector: 'app-pie-chart',
  template: `<div #pieChart style="height: 250px;"></div>`
})
export class PieChartComponent implements AfterViewInit{

  @ViewChild('pieChart') pieChart: ElementRef;
  @Input() data: any;
  @Input() title: any;

  drawChart = () => {
    const data = google.visualization.arrayToDataTable(this.data);

    const options = {
      title: this.title,
      legend: {position: 'right'},
      colors: ['#7cd6fd', '#743ee2', '#5e64ff',  '#ff5858', '#ffa00a'],
    };

    const chart = new google.visualization.PieChart(this.pieChart.nativeElement);

    chart.draw(data, options);
  }

  ngAfterViewInit() {
    google.charts.load('current', { packages: ['corechart'] });
    google.charts.setOnLoadCallback(this.drawChart);
  }
}