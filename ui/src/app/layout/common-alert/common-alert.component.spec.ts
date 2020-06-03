import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { CommonAlertComponent } from './common-alert.component';

describe('CommonAlertComponent', () => {
  let component: CommonAlertComponent;
  let fixture: ComponentFixture<CommonAlertComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ CommonAlertComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(CommonAlertComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
