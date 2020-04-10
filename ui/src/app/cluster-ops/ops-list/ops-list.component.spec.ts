import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { OpsListComponent } from './ops-list.component';

describe('OpsListComponent', () => {
  let component: OpsListComponent;
  let fixture: ComponentFixture<OpsListComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ OpsListComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(OpsListComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
