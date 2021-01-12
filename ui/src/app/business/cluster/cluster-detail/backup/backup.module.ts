import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BackupListComponent } from './backup-list/backup-list.component';
import { BackupStrategyComponent } from './backup-strategy/backup-strategy.component';
import {CoreModule} from '../../../../core/core.module';
import { BackupLogComponent } from './backup-log/backup-log.component';



@NgModule({
    declarations: [BackupListComponent, BackupStrategyComponent, BackupLogComponent],
    exports: [
        BackupListComponent,
        BackupStrategyComponent,
        BackupLogComponent
    ],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class BackupModule { }
