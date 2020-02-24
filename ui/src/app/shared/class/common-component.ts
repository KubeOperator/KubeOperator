import {OnDestroy, OnInit} from '@angular/core';
import {ItemChangeService} from '../../base/header/components/item-change/item-change.service';
import {Profile} from '../session-user';
import {Subscription} from 'rxjs';

export abstract class CommonComponent implements OnInit, OnDestroy {
  public profile: Profile;
  public sub: Subscription;

  protected constructor(public itemChangeService: ItemChangeService) {
  }

  ngOnDestroy(): void {
    this.sub.unsubscribe();
  }

  ngOnInit(): void {
    this.sub = this.itemChangeService.$noticeChannel.subscribe(data => {
      this.profile = data;
    });
  }
}
