import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { PlanDeleteComponent } from './plan-delete.component';

describe('PlanDeleteComponent', () => {
  let component: PlanDeleteComponent;
  let fixture: ComponentFixture<PlanDeleteComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ PlanDeleteComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(PlanDeleteComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
