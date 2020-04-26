import { Component } from '@angular/core';
import { DataService } from './providers/data.service';
import { Data } from './interface/data.interface';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, combineLatest, EMPTY, of } from 'rxjs';
import { map, catchError, tap } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {

  private filterDependenciesSelectedSubject = new BehaviorSubject<string>('*');
  filterDependenciesSelectedAction$ = this.filterDependenciesSelectedSubject.asObservable();

  private dependenciesByVersionsSubject = new BehaviorSubject<any>({});
  dependenciesByVersions$ = this.dependenciesByVersionsSubject.asObservable();

  public data$ = combineLatest([
    this.dataService.cacheData$,
    this.filterDependenciesSelectedAction$,
  ]).pipe(
    map(([data, filterDependencies]) => {
      const dependenciesFiltered = Object.keys(data.dependenciesByVersions)
        .filter((key) => {
          return filterDependencies === '*'
            ? true
            : data.dependenciesByVersions[key].type === filterDependencies;
        })
        .reduce((obj, key) => {
          obj[key] = data.dependenciesByVersions[key];
          return obj;
        }, {});

      this.dependenciesByVersionsSubject.next(dependenciesFiltered);
      return data;
    })
  );
  public serverIsOn$ = this.dataService.serverIsOn$;

  objectKeys = Object.keys;
  showLoading = false;

  constructor(
    private dataService: DataService,
    private translate: TranslateService
  ) {
    const browserLang = this.translate.getBrowserLang();
    translate.setDefaultLang(browserLang.match(/en|pt/) ? browserLang : 'en');
  }

  refresh() {
    this.data$ = this.dataService.data$;
  }

  filterDependencies(type: string) {
    this.filterDependenciesSelectedSubject.next(type);
  }
}
