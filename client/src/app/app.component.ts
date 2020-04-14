import { Component, OnInit, OnDestroy } from '@angular/core';
import { DataService } from './providers/data.service';
import { Data } from './interface/data.interface';
import { TranslateService } from '@ngx-translate/core';
import { Subject } from 'rxjs';
import { takeUntil, take } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {
    private destroy$ = new Subject<boolean>();

    data: Data;
    objectKeys = Object.keys;
    showLoading = false;
    showRefreshButton = false;

    constructor(
        private dataService: DataService,
        private translate: TranslateService) {
            const browserLang = this.translate.getBrowserLang();
            translate.setDefaultLang(browserLang.match(/en|pt/) ? browserLang : 'en');
        }

    ngOnInit() {
        this.getCache();
        this.checkServerIsOn();
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

    checkServerIsOn() {
        this.dataService.ping()
        .pipe(
            take(1),
            takeUntil(this.destroy$)
        )
        .subscribe(value => {
          if (value.ok === true) {
            this.showRefreshButton = true;
          }
        });
    }

    ngOnDestroy() {
        this.destroy$.next(true);
        this.destroy$.unsubscribe();
    }
}
