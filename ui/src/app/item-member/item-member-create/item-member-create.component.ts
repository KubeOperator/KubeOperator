import {Component, OnInit} from '@angular/core';
import {ItemMemberService} from '../item-member.service';
import {Profile} from '../../shared/session-user';


@Component({
  selector: 'app-item-member-create',
  templateUrl: './item-member-create.component.html',
  styleUrls: ['./item-member-create.component.css']
})


export class ItemMemberCreateComponent implements OnInit {
  opened = true;
  ps: Profile[] = [];
  options = [];
  managers = [];
  viewers = [];
  ops: any = {
    multiple: true,
    placeholder: '选择用户',
    escapeMarkup: function (markup) {
      return markup;
    },
    templateSelection: (data) => {
      return `<span class="label label-purple select2-selection__choice__remove">${data['text']}</span>`;
    }
  };

  userFilter() {
    const all = this.managers.concat(this.viewers);
    const remove = [];
    this.options = this.toOptions();
    this.options.forEach(o => {
      all.forEach(a => {
        if (o['value'] === a['value']) {
          remove.push(o);
        }
      });
    });
    remove.forEach(r => {
      r['disabled'] = true;
    });
    console.log(this.options);
  }

  constructor(private itemMemberService: ItemMemberService) {
  }

  ngOnInit() {
    this.listUser();
  }

  listUser() {
    this.itemMemberService.getProfiles().subscribe(data => {
      data.forEach(p => {
        this.ps.push(p);
      });
      this.options = this.toOptions();
      console.log(this.options);
    });
  }

  toOptions(): any[] {
    const options = [];
    this.ps.forEach(p => {
      options.push({'id': p.id, 'text': p.user.username, 'value': p.id});
    });
    return options;
  }
}
