import {Component, OnInit, ViewChild} from '@angular/core';
import {CredentialCreateComponent} from '../credential/credential-create/credential-create.component';
import {CredentialListComponent} from '../credential/credential-list/credential-list.component';
import {RegionListComponent} from './region-list/region-list.component';
import {RegionCreateComponent} from './region-create/region-create.component';

@Component({
  selector: 'app-region',
  templateUrl: './region.component.html',
  styleUrls: ['./region.component.css']
})
export class RegionComponent implements OnInit {

  @ViewChild(RegionCreateComponent)
  creation: RegionCreateComponent;

  @ViewChild(RegionListComponent)
  listRegion: RegionListComponent;

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
    this.listRegion.refresh();
  }

}
