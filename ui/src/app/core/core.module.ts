import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {ClarityModule} from '@clr/angular';
import {FormsModule} from '@angular/forms';
import {TranslateModule} from '@ngx-translate/core';
import {LSelect2Module} from 'ngx-select2';
import {NgxEchartsModule} from 'ngx-echarts';


@NgModule({
    declarations: [],
    imports: [
        CommonModule,
        FormsModule,
        ClarityModule,
        TranslateModule,
        LSelect2Module,
        NgxEchartsModule,
    ],
    exports: [
        CommonModule,
        FormsModule,
        ClarityModule,
        TranslateModule,
        LSelect2Module,
    ]
})
export class CoreModule {
}
