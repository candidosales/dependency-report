import { Component, OnInit } from '@angular/core';
import { DataService } from './providers/data.service';
import { Data } from './interface/data.interface';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {

    data: Data;
    objectKeys = Object.keys;

    constructor(
        private dataService: DataService,
        private translate: TranslateService) {
            const browserLang = this.translate.getBrowserLang();
            translate.setDefaultLang(browserLang.match(/en|pt/) ? browserLang : 'en');
        }

    ngOnInit() {
        this.dataService.getDataInServer()
        .subscribe(value => {
            value.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
            value.graphData.componentsByFilters?.unshift(['Filter', 'Version']);

            this.data = value;
        });
    }

    getVersionByFilter(filter: string): string {
        const values = filter.split('_');
        return values[1];
    }

    getIconByFilter(filter: string): string {
        if (filter.includes('angular')) {
            return 'https://coryrylan.com/assets/images/posts/types/angular.svg';
        }
        return '';
    }

    refresh() {
        this.dataService.generateReport().subscribe(httpResponse => {
            console.log('httpResponse', httpResponse);
        });
    }
}
