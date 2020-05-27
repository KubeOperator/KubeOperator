import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {SettingComponent} from './setting.component';
import {RouterModule} from '@angular/router';
import {TranslateModule} from '@ngx-translate/core';


@NgModule({
    declarations: [SettingComponent],
    imports: [
        CommonModule,
        RouterModule,
        TranslateModule
    ]
})

export class SettingModule {

}
