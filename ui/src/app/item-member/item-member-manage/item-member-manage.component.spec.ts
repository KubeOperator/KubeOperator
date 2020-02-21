import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { ItemMemberManageComponent } from './item-member-manage.component';

describe('ItemMemberManageComponent', () => {
  let component: ItemMemberManageComponent;
  let fixture: ComponentFixture<ItemMemberManageComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ ItemMemberManageComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ItemMemberManageComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
