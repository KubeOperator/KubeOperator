import {Component, EventEmitter, OnInit, Output, ViewChild} from '@angular/core';
import {User} from '../../../../user/user';
import {UserService} from '../../../../user/user.service';
import {SessionService} from '../../../../shared/session.service';
import {NgForm} from '@angular/forms';
import {Profile} from '../../../../shared/session-user';

@Component({
  selector: 'app-item-change',
  templateUrl: './item-change.component.html',
  styleUrls: ['./item-change.component.css']
})
export class ItemChangeComponent implements OnInit {
  profile: Profile;
  isSubmitGoing = false;
  opened = false;
  loading = false;
  @ViewChild('itemForm', {static: true})
  itemForm: NgForm;
  @Output() changeItem = new EventEmitter();

  constructor(private userService: UserService, private sessionService: SessionService) {
  }

  ngOnInit() {
  }

  open() {
    this.itemForm.resetForm();
    this.opened = true;
    this.profile = this.sessionService.getCacheProfile();
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    this.sessionService.changeItem(this.profile.current_item).subscribe(() => {
      this.sessionService.cacheProfile(this.profile);
      this.changeItem.emit();
      this.opened = false;
      this.isSubmitGoing = false;
    });
  }

  onCancel() {
    this.opened = false;
  }

}
