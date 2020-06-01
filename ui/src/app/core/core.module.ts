import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {AppRoutingModule} from '../app-routing.module';
import {ClarityModule} from '@clr/angular';
import {FormsModule} from '@angular/forms';
import {TranslateModule} from '@ngx-translate/core';
import {LSelect2Module} from 'ngx-select2';


@NgModule({
    declarations: [],
    imports: [
        CommonModule,
        FormsModule,
        ClarityModule,
        TranslateModule,
        LSelect2Module,
    ],
    exports: [
        CommonModule,
        FormsModule,
        ClarityModule,
        TranslateModule,
        LSelect2Module
    ]
})
export class CoreModule {
}
