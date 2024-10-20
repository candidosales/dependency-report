import { AfterViewInit, Component, ElementRef, ViewChild, HostListener, OnChanges, SimpleChanges, ChangeDetectionStrategy, input } from '@angular/core';
declare var google: any;

@Component({
    selector: 'app-pie-chart',
    templateUrl: './pie-chart.component.html',
    changeDetection: ChangeDetectionStrategy.OnPush,
    standalone: true
})
export class PieChartComponent implements AfterViewInit, OnChanges {

  @ViewChild('pieChart') pieChart: ElementRef;
  data = input([['empty', 0]]);
  title = input('');
  height = input('250px');
  width = input('100%');

  @HostListener('window:resize', ['$event'])
  onWindowResize(event: any) {
    this.drawChart();
  }

  drawChart = () => {
    if (this.data() && google.visualization !== undefined && this.pieChart !== undefined) {
      try {
        const data = google.visualization.arrayToDataTable(this.data());
        const options = {
          title: this.title(),
          legend: {position: 'right'},
          colors: ['#D53F8C', '#805AD5', '#5A67D8', '#3182CE', '#319795',  '#38A169', '#D69E2E'],
        };

        const chart = new google.visualization.PieChart(this.pieChart.nativeElement);
        chart.draw(data, options);
      } catch (err) {
        console.error(err);
      }
    }
  }

  ngAfterViewInit() {
    google.charts.load('current', { packages: ['corechart'] });
    google.charts.setOnLoadCallback(this.drawChart);
  }

  ngOnChanges(changes: SimpleChanges): void {
    this.drawChart();
  }
}
