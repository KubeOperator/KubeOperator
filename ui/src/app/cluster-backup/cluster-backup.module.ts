import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ClusterBackupListComponent } from './cluster-backup-list/cluster-backup-list.component';
import {CoreModule} from '../core/core.module';
import { ClusterBackupStrategyComponent } from './cluster-backup-strategy/cluster-backup-strategy.component';
import {SharedModule} from '../shared/shared.module';


@NgModule({
  declarations: [ClusterBackupListComponent, ClusterBackupStrategyComponent],
  exports: [
    ClusterBackupListComponent,
    ClusterBackupStrategyComponent
  ],
  imports: [
    CommonModule,
    CoreModule,

    SharedModule,
  ]
})
export class ClusterBackupModule { }
