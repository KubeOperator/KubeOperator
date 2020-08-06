import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { BackupListComponent } from './backup-list/backup-list.component';
import { BackupStrategyComponent } from './backup-strategy/backup-strategy.component';
import {CoreModule} from '../../../../core/core.module';



@NgModule({
    declarations: [BackupListComponent, BackupStrategyComponent],
    exports: [
        BackupListComponent,
        BackupStrategyComponent
    ],
    imports: [
        CommonModule,
        CoreModule
    ]
})
export class BackupModule { }
