import { async, ComponentFixture, TestBed } from '@angular/core/testing';

import { DeployPlanComponent } from './deploy-plan.component';

describe('DeployPlanComponent', () => {
  let component: DeployPlanComponent;
  let fixture: ComponentFixture<DeployPlanComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      declarations: [ DeployPlanComponent ]
    })
    .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(DeployPlanComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
