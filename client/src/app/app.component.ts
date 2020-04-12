import { Component, OnInit } from '@angular/core';
import { DataService } from './providers/data.service';
import { Data } from './interface/data.interface';
import { TranslateService } from '@ngx-translate/core';
import { Subject } from 'rxjs';
import { takeUntil } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit {
    private destroy$ = new Subject<boolean>();

    data: Data;
    objectKeys = Object.keys;
    showLoading = false;

    constructor(
        private dataService: DataService,
        private translate: TranslateService) {
            const browserLang = this.translate.getBrowserLang();
            translate.setDefaultLang(browserLang.match(/en|pt/) ? browserLang : 'en');
        }

    ngOnInit() {
        this.getCache();
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
        this.showLoading = true;
        this.dataService.generateReport()
        .pipe(takeUntil(this.destroy$))
        .subscribe(value => {
            this.showLoading = false;
            this.prepareData(value);
        }, error => {
            this.getCache();
            this.showLoading = false;
        });
    }

    getCache() {
        this.showLoading = true;
        this.dataService.getCacheData()
        .pipe(takeUntil(this.destroy$))
        .subscribe(value => {
            this.showLoading = false;
            this.prepareData(value);
        }, error => {
            this.showLoading = false;
        });
    }

    prepareData(value: Data) {
        value.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
        value.graphData.componentsByFilters?.unshift(['Filter', 'Version']);
        this.data = value;
    }
}
