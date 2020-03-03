import {Component, EventEmitter, Input, OnInit, Output} from '@angular/core';
import {ItemMemberService} from '../item-member.service';
import {Profile} from '../../shared/session-user';
import {ActivatedRoute} from '@angular/router';
import {SessionService} from '../../shared/session.service';

const ROLE_NAME_MANAGER = 'MANAGER';
const ROLE_NAME_VIEWER = 'VIEWER';


@Component({
  selector: 'app-item-member-create',
  templateUrl: './item-member-create.component.html',
  styleUrls: ['./item-member-create.component.css']
})


export class ItemMemberCreateComponent implements OnInit {
  opened = false;
  loading = true;
  isSubmitGoing = false;
  @Output() create = new EventEmitter<boolean>();
  currentItem;
  ps: Profile[] = [];
  options = [];
  managers = [];
  viewers = [];
  itemName: string;
  profile: Profile;
  ops: any = {
    multiple: true,
    placeholder: '选择用户',
    escapeMarkup: function (markup) {
      return markup;
    },
    templateSelection: (data) => {
      return `<span class="label label-blue select2-selection__choice__remove">${data['text']}</span>`;
    },
  };

  userFilter() {
    const all = this.managers.concat(this.viewers);
    const remove = [];
    this.options = this.toOptions(this.ps);
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
  }

  constructor(private itemMemberService: ItemMemberService, private route: ActivatedRoute, private session: SessionService) {
  }

  ngOnInit() {
    this.profile = this.session.getCacheProfile();
    this.currentItem = this.route.snapshot.queryParams['name'];
  }

  onCancel() {
    this.opened = false;
  }

  onSubmit() {
    if (this.isSubmitGoing) {
      return;
    }
    this.isSubmitGoing = true;
    const submitData = {'role_map': {}, 'profiles': []};
    const roleMap = {};
    this.viewers.forEach(v => {
      roleMap[v['value']] = ROLE_NAME_VIEWER;
    });
    this.managers.forEach(v => {
      roleMap[v['value']] = ROLE_NAME_MANAGER;
    });
    this.viewers.concat(this.managers).forEach(p => {
      submitData.profiles.push(p['value']);
      submitData.role_map = roleMap;
    });
    this.itemMemberService.setItemProfiles(submitData, this.currentItem).subscribe(() => {
      this.isSubmitGoing = false;
      this.create.emit(true);
      this.opened = false;
    });
  }

  open(profiles: Profile[]) {
    this.clear();
    this.opened = true;
    this.itemName = this.currentItem;
    this.itemMemberService.getProfiles().subscribe(data => {
      data.filter((p) => {
        return !p.user.is_superuser;
      }).forEach(p => {
        this.ps.push(p);
      });
      this.viewers = this.getOptionsByRole(profiles, ROLE_NAME_VIEWER);
      this.managers = this.getOptionsByRole(profiles, ROLE_NAME_MANAGER);
      this.userFilter();
    });
  }

  clear() {
    this.options = [];
    this.managers = [];
    this.viewers = [];
    this.ps = [];
  }

  toOptions(profiles: Profile[]): any[] {
    const options = [];
    profiles.forEach(p => {
      options.push({'id': p.id, 'text': p.user.username, 'value': p.id});
    });
    return options;
  }

  getOptionsByRole(profiles: Profile[], roleName: string): any[] {
    const m = profiles.filter((p) => {
      return this.formatRole(p) === roleName;
    });
    return this.toOptions(m);
  }

  formatRole(p: Profile) {
    return p.item_role_mappings.find(mp => {
      return mp.item_name === this.itemName;
    })['role'];
  }

}
