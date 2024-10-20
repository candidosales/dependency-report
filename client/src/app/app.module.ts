import { AppComponent } from './app.component';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { BrowserModule } from '@angular/platform-browser';
import { DataService } from './providers/data.service';
import { DependencyComponent } from './dependency/dependency.component';
import { HttpClient, HttpClientModule } from '@angular/common/http';
import { IconDesignPipe } from './pipes/icon-design.pipe';
import { IconLibraryPipe } from './pipes/icon-library.pipe';
import { MatIconModule } from '@angular/material/icon';
import { MatLegacyProgressSpinnerModule as MatProgressSpinnerModule } from '@angular/material/legacy-progress-spinner';
import { MatLegacyTooltipModule as MatTooltipModule } from '@angular/material/legacy-tooltip';
import { NgModule } from '@angular/core';
import { NotificationsDialogComponent } from './repository/notifications-dialog/notifications-dialog.component';
import { PieChartComponent } from './google-charts/pie-chart.component';
import { RepositoryComponent } from './repository/repository.component';
import { TranslateHttpLoader } from '@ngx-translate/http-loader';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { TreemapChartComponent } from './google-charts/treemap-chart.component';
import { UpdatesDialogComponent } from './repository/updates-dialog/updates-dialog.component';
import { VersionPipe } from './pipes/version.pipe';
import { WithLoadingPipe } from './pipes/with-loading.pipe';
import { MatDialogModule } from '@angular/material/dialog';
import { MatButtonModule } from '@angular/material/button';
import { MatSnackBarModule } from '@angular/material/snack-bar';

export function HttpLoaderFactory(http: HttpClient) {
  return new TranslateHttpLoader(http);
}

@NgModule({
  declarations: [
    AppComponent,
    PieChartComponent,
    TreemapChartComponent,
    DependencyComponent,
    RepositoryComponent,
    IconDesignPipe,
    IconLibraryPipe,
    VersionPipe,
    WithLoadingPipe,
    NotificationsDialogComponent,
    UpdatesDialogComponent
  ],
  imports: [
    BrowserModule,
    HttpClientModule,
    BrowserAnimationsModule,
    MatTooltipModule,
    MatButtonModule,
    MatIconModule,
    MatProgressSpinnerModule,
    MatDialogModule,
    MatSnackBarModule,
    TranslateModule.forRoot({
        loader: {
            provide: TranslateLoader,
            useFactory: HttpLoaderFactory,
            deps: [HttpClient]
        }
    })
  ],
  providers: [DataService],
  bootstrap: [AppComponent]
})
export class AppModule { }
