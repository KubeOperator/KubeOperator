import {NgModule} from '@angular/core';
import {HeaderComponent} from './header/header.component';
import {FooterComponent} from './footer/footer.component';
import {LayoutComponent} from './layout.component';
import {NavigationComponent} from './navigation/navigation.component';
import {AppAlertComponent} from './app-alert/app-alert.component';
import {CoreModule} from '../core/core.module';
import {RouterModule} from '@angular/router';
import {SharedModule} from '../shared/shared.module';
import {CommonAlertComponent} from './common-alert/common-alert.component';
import {PasswordComponent} from './header/password/password.component';
import { AboutComponent } from './header/about/about.component';

@NgModule({
    declarations: [HeaderComponent, FooterComponent, LayoutComponent,
        NavigationComponent, AppAlertComponent, CommonAlertComponent, PasswordComponent, AboutComponent],
    exports: [
        LayoutComponent,
        CommonAlertComponent,
    ],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
    ]
})
export class LayoutModule {
}
