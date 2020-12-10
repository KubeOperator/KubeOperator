import {BrowserModule} from '@angular/platform-browser';
import {ErrorHandler, NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {TranslateLoader, TranslateModule} from '@ngx-translate/core';
import {LoginModule} from './login/login.module';
import {AppRoutingModule} from './app-routing.module';
import {LayoutModule} from './layout/layout.module';
import {HTTP_INTERCEPTORS, HttpClient, HttpClientModule} from '@angular/common/http';
import {TranslateHttpLoader} from '@ngx-translate/http-loader';
import {BusinessModule} from './business/business.module';
import {ClrCommonFormsModule, ClrDatagridModule, ClrIconModule, ClrModalModule} from '@clr/angular';
import {CoreModule} from './core/core.module';
import {SessionInterceptor} from './shared/auth/session-interceptor';
import {AppGlobalErrorHandler} from './shared/handler/app-global-error-handler';
import {NgCircleProgressModule} from 'ng-circle-progress';
import {CodemirrorModule} from 'ng2-codemirror';


export function HttpLoaderFactory(httpClient: HttpClient) {
    return new TranslateHttpLoader(httpClient, './assets/i18n/', '.json?timestamp=' + Math.random());
}

@NgModule({
    declarations: [
        AppComponent
    ],
    imports: [
        BrowserModule,
        BrowserAnimationsModule,
        AppRoutingModule,
        LayoutModule,
        BusinessModule,
        LoginModule,
        HttpClientModule,
        TranslateModule.forRoot({
            loader: {
                provide: TranslateLoader,
                useFactory: HttpLoaderFactory,
                deps: [HttpClient]
            }
        }),
        ClrDatagridModule,
        ClrIconModule,
        ClrModalModule,
        ClrCommonFormsModule,
        CoreModule,
        CodemirrorModule,
        NgCircleProgressModule.forRoot({
            radius: 100,
            outerStrokeWidth: 16,
            innerStrokeWidth: 8,
            animationDuration: 300,
        })
    ],
    providers: [
        {
            provide: HTTP_INTERCEPTORS,
            useClass: SessionInterceptor,
            multi: true,
        },
        {provide: ErrorHandler, useClass: AppGlobalErrorHandler},
    ],
    bootstrap: [AppComponent],

})
export class AppModule {
}
