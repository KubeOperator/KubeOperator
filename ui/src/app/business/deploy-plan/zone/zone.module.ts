import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ZoneCreateComponent } from './zone-create/zone-create.component';
import { ZoneDeleteComponent } from './zone-delete/zone-delete.component';
import { ZoneUpdateComponent } from './zone-update/zone-update.component';
import { ZoneListComponent } from './zone-list/zone-list.component';



@NgModule({
  declarations: [ZoneCreateComponent, ZoneDeleteComponent, ZoneUpdateComponent, ZoneListComponent],
  imports: [
    CommonModule
  ]
})
export class ZoneModule { }
