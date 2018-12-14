import {Component, OnInit, ViewChild} from '@angular/core';
import {UserCreateComponent} from './user-create/user-create.component';
import {UserListComponent} from './user-list/user-list.component';

@Component({
  selector: 'app-user',
  templateUrl: './user.component.html',
  styleUrls: ['./user.component.css']
})
export class UserComponent implements OnInit {

  @ViewChild(UserCreateComponent)
  creationUser: UserCreateComponent;
  @ViewChild(UserListComponent)
  listUser: UserListComponent;

  constructor() {
  }

  openModal() {
    this.creationUser.newUser();
  }

  createUser(created: boolean) {
    if (created) {
      this.refresh();
    }
  }

  refresh() {
    this.listUser.listUser();
  }

  ngOnInit() {
  }

}
