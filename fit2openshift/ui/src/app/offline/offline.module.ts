import {NgModule} from '@angular/core';
import {CommonModule} from '@angular/common';
import {OfflineComponent} from './offline.component';
import {OfflineService} from './offline.service';
import {SharedModule} from '../shared/shared.module';
import {OfflineListComponent} from './offline-list/offline-list.component';

@NgModule({
  declarations: [OfflineComponent, OfflineListComponent],
  imports: [
    CommonModule,
    SharedModule
  ], providers: [OfflineService]
})
export class OfflineModule {
}
