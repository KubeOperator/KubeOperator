import {Component, OnInit, ViewChild} from '@angular/core';
import {ZoneCreateComponent} from './zone-create/zone-create.component';
import {ZoneListComponent} from './zone-list/zone-list.component';

@Component({
  selector: 'app-zone',
  templateUrl: './zone.component.html',
  styleUrls: ['./zone.component.css']
})
export class ZoneComponent implements OnInit {

  @ViewChild(ZoneCreateComponent, { static: true })
  creation: ZoneCreateComponent;

  @ViewChild(ZoneListComponent, { static: true })
  listZone: ZoneListComponent;

  constructor() {
  }

  ngOnInit() {
  }

  openModal() {
    this.creation.newItem();
  }

  create(created: boolean) {
    if (created) {
      this.refresh();
    }
  }

  refresh() {
    this.listZone.refresh();
  }

}
