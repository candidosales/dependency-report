import { Component, inject } from '@angular/core';
import { DataService } from './providers/data.service';
import { TranslateService } from '@ngx-translate/core';
import { BehaviorSubject, combineLatest } from 'rxjs';
import { map, catchError, tap, filter } from 'rxjs/operators';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.scss'],
})
export class AppComponent {
  private dataService = inject(DataService);
  private translate = inject(TranslateService);


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
  public disableRefreshButton$ = this.dataService.serverIsOn$.pipe(
    map(value => {
      return !value;
    })
  );

  objectKeys = Object.keys;
  showLoading = false;

  constructor() {
    const translate = this.translate;

    const browserLang = this.translate.getBrowserLang();
    translate.setDefaultLang(browserLang.match(/en|pt/) ? browserLang : 'en');
  }

  refresh() {
    this.showLoading = true;
    this.data$ = this.dataService.data$.pipe(
      filter((value) => !!value),
      tap(value => this.showLoading = false),
      catchError(err => {
        console.error(err);
        return this.dataService.cacheData$;
      })
    );
  }

  filterDependencies(type: string) {
    this.filterDependenciesSelectedSubject.next(type);
  }
}
