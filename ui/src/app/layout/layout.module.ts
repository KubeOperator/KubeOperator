import {NgModule} from '@angular/core';
import {HeaderComponent} from './header/header.component';
import {FooterComponent} from './footer/footer.component';
import {LayoutComponent} from './layout.component';
import {NavigationComponent} from './navigation/navigation.component';
import {AppAlertComponent} from './app-alert/app-alert.component';
import {CoreModule} from '../core/core.module';
import {RouterModule} from '@angular/router';
import {SharedModule} from '../shared/shared.module';


@NgModule({
    declarations: [HeaderComponent, FooterComponent, LayoutComponent, NavigationComponent, AppAlertComponent],
    exports: [
        LayoutComponent
    ],
    imports: [
        CoreModule,
        RouterModule,
        SharedModule,
    ]
})
export class LayoutModule {
}
