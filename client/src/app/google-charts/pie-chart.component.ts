import { AfterViewInit, Component, ElementRef, ViewChild, Input, HostListener, OnChanges, SimpleChanges, ChangeDetectionStrategy } from '@angular/core';
declare var google: any;

@Component({
  selector: 'app-pie-chart',
  templateUrl: './pie-chart.component.html',
  changeDetection: ChangeDetectionStrategy.OnPush
})
export class PieChartComponent implements AfterViewInit, OnChanges {

  @ViewChild('pieChart') pieChart: ElementRef;
  @Input() data = [[ 'empty', 0 ]];
  @Input() title = '';
  @Input() height = '250px';
  @Input() width = '100%';

  @HostListener('window:resize', ['$event'])
  onWindowResize(event: any) {
    this.drawChart();
  }

  drawChart = () => {
    if (this.data && google.visualization !== undefined) {
      try {
        const data = google.visualization.arrayToDataTable(this.data);
        const options = {
          title: this.title,
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
