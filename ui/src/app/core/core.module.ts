import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {AppRoutingModule} from '../app-routing.module';
import {ClarityModule} from '@clr/angular';
import {FormsModule} from '@angular/forms';
import {TranslateModule} from '@ngx-translate/core';


@NgModule({
    declarations: [],
    imports: [
        CommonModule,
        FormsModule,
        ClarityModule,
        TranslateModule
    ],
    exports: [
        CommonModule,
        FormsModule,
        ClarityModule,
        TranslateModule
    ]
})
export class CoreModule {
}
