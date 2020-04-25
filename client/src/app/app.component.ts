import { Component, OnInit, OnDestroy } from '@angular/core';
import { DataService } from './providers/data.service';
import { Data } from './interface/data.interface';
import { TranslateService } from '@ngx-translate/core';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss']
})
export class AppComponent implements OnInit, OnDestroy {
    public data$ = this.dataService.cacheData$;
    public serverIsOn$ = this.dataService.serverIsOn$;

    objectKeys = Object.keys;
    showLoading = false;

    constructor(
        private dataService: DataService,
        private translate: TranslateService) {
            const browserLang = this.translate.getBrowserLang();
            translate.setDefaultLang(browserLang.match(/en|pt/) ? browserLang : 'en');
        }

    ngOnInit() {
    }

    refresh() {
        this.data$ = this.dataService.data$;
    }

    prepareData(value: Data) {
        value.graphData.projectsByFilters?.unshift(['Filter', 'Version']);
        value.graphData.componentsByFilters?.unshift(['Filter', 'Version']);
    }

    ngOnDestroy() {
        this.destroy$.next(true);
        this.destroy$.unsubscribe();
    }
}
