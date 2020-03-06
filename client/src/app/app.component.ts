import { Component } from '@angular/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent {
  title = 'dependency-report';

  dataSet = {
      labels: ['@angular/7', '@angular/8', '@angular/9'],
      datasets: [
        {
          values: [5, 34, 10]
        },
      ]
  };

  projects = [
              {
                  'name' : 'listing-builder-client',
                  'version' : '3.4.5',
                  'filter': '@angular/core_8'
              },
              {
                  'name' : 'business-center-client',
                  'version' : '3.4.5',
                  'filter': '@angular/core_8'
              },
              {
                  'name' : 'reputation-client',
                  'version' : '3.4.5',
                  'filter': '@angular/core_9'
              },
              {
                  'name' : 'customer-voice-client',
                  'version' : '3.4.5',
                  'filter': '@angular/core_9'
              }
          ];
          components = [
            {
                'name' : 'uikit',
                'version' : '3.4.5',
                'filter': '@angular/core_8'
            },
            {
                'name' : 'forms',
                'version' : '8.0.15',
                'filter': '@angular/core_8'
            },
            {
                'name' : 'core',
                'version' : '9.0.5',
                'filter': '@angular/core_9'
            },
        ];
}
