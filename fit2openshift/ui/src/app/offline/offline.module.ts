import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {OfflineComponent} from './offline.component';
import {OfflineService} from './offline.service';
import {SharedModule} from '../shared/shared.module';
import {OfflineListComponent} from './offline-list/offline-list.component';
import {CoreModule} from '../core/core.module';
import {TipModule} from '../tip/tip.module';

@NgModule({
  declarations: [OfflineComponent, OfflineListComponent],
  imports: [
    CommonModule,
    SharedModule,
    CoreModule,
    TipModule
  ], providers: [OfflineService]
})
export class OfflineModule {
}
