import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {AppComponent} from './app.component';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {TranslateLoader, TranslateModule} from '@ngx-translate/core';
import {LoginModule} from './login/login.module';
import {AppRoutingModule} from './app-routing.module';
import {LayoutModule} from './layout/layout.module';
import {HttpClient, HttpClientModule} from '@angular/common/http';
import {TranslateHttpLoader} from '@ngx-translate/http-loader';
import {BusinessModule} from './business/business.module';
import { CredentialComponent } from './business/setting/credential/credential.component';
import { CredentialListComponent } from './business/setting/credential/credential-list/credential-list.component';
import { CredentialCreateComponent } from './business/setting/credential/credential-create/credential-create.component';
import {ClrCommonFormsModule, ClrDatagridModule, ClrIconModule, ClrModalModule} from '@clr/angular';
import {CoreModule} from './core/core.module';


export function HttpLoaderFactory(httpClient: HttpClient) {
    return new TranslateHttpLoader(httpClient);
}

@NgModule({
    declarations: [
        AppComponent,
        CredentialComponent,
        CredentialListComponent,
        CredentialCreateComponent,
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
        CoreModule
    ],
    providers: [],
    bootstrap: [AppComponent],

})
export class AppModule {
}
