import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { UserDeleteComponent } from './user-delete.component';

describe('UserDeleteComponent', () => {
  let component: UserDeleteComponent;
  let fixture: ComponentFixture<UserDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ UserDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(UserDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
