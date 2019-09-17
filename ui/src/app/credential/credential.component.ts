import {Component, OnInit, ViewChild} from '@angular/core';
import {CredentialCreateComponent} from './credential-create/credential-create.component';
import {CredentialListComponent} from './credential-list/credential-list.component';

@Component({
  selector: 'app-credential',
  templateUrl: './credential.component.html',
  styleUrls: ['./credential.component.css']
})
export class CredentialComponent implements OnInit {

  @ViewChild(CredentialCreateComponent, { static: true })
  creation: CredentialCreateComponent;

  @ViewChild(CredentialListComponent, { static: true })
  listHost: CredentialListComponent;

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
    this.listHost.refresh();
  }
}
